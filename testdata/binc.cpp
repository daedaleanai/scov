#include <algorithm>
#include <cassert>
#include <cstdio>
#include <string>

#include <llvm/IR/LLVMContext.h>
#include <llvm/IR/Module.h>
#include <llvm/MC/SubtargetFeature.h>
#include <llvm/Support/CommandLine.h>
#include <llvm/Support/FileSystem.h>
#if LLVM_CONFIG >= 700
#include <llvm/Support/InitLLVM.h>
#endif
#include <llvm/Support/Path.h>
#include <llvm/Support/TargetRegistry.h>
#include <llvm/Support/TargetSelect.h>
#include <llvm/Target/TargetMachine.h>

#include "binc.h"

namespace cl = llvm::cl;

using std::string;

enum CodeGenFileType {
	CGFT_AssemblyFile = llvm::TargetMachine::CGFT_AssemblyFile,
	CGFT_ObjectFile = llvm::TargetMachine::CGFT_ObjectFile,
	CGFT_Null = llvm::TargetMachine::CGFT_Null,
	CGFT_IR
};

cl::OptionCategory category( "binc" );
cl::opt<string> OutputFilename( "o", cl::desc( "Specify the output filename.  If -, use stdout." ),
								cl::value_desc( "filename" ), cl::cat( category ) );
cl::opt<string> InputFilename( "c", cl::desc( "Specify the input filename.  If -, use stdin." ),
							   cl::value_desc( "filename" ), cl::Required, cl::init( "-" ), cl::cat( category ) );
cl::opt<string> VariableName( "n", cl::desc( "Specify linkage name for binary asset." ), cl::value_desc( "identifier" ),
							  cl::cat( category ) );
cl::opt<bool> NullTerminate( "z", cl::desc( "Add a null terminator to the binary asset." ), cl::cat( category ) );
cl::opt<bool> DebugInfo( "g", cl::desc( "Turn on debugging information" ), cl::init( false ), cl::cat( category ) );
cl::opt<std::string> MArch( "march", cl::desc( "Architecture to generate code for (see --version)" ),
							cl::cat( category ) );
cl::opt<bool> DataSections( "data-sections", cl::desc( "Emit data into separate sections" ), cl::init( false ),
							cl::cat( category ) );
cl::opt<std::string> MCPU( "mcpu", cl::desc( "Target a specific cpu type (-mcpu=help for details)" ),
						   cl::value_desc( "cpu-name" ), cl::init( "" ), cl::cat( category ) );
cl::list<std::string> MAttrs( "mattr", cl::CommaSeparated,
							  cl::desc( "Target specific attributes (-mattr=help for details)" ),
							  cl::value_desc( "a1,+a2,-a3,..." ), cl::cat( category ) );
cl::opt<std::string> TargetTriple( "mtriple", cl::desc( "Override target triple for module" ), cl::cat( category ) );
cl::opt<llvm::CodeModel::Model>
	CMModel( "code-model", cl::desc( "Choose code model" ),
			 cl::values( clEnumValN( llvm::CodeModel::Small, "small", "Small code model" ),
						 clEnumValN( llvm::CodeModel::Kernel, "kernel", "Kernel code model" ),
						 clEnumValN( llvm::CodeModel::Medium, "medium", "Medium code model" ),
						 clEnumValN( llvm::CodeModel::Large, "large", "Large code model" ) ),
			 cl::cat( category ) );
cl::opt<llvm::Reloc::Model> RelocModel(
	"relocation-model", cl::desc( "Choose relocation model" ),
	cl::values( clEnumValN( llvm::Reloc::Static, "static", "Non-relocatable code" ),
				clEnumValN( llvm::Reloc::PIC_, "pic", "Fully relocatable, position independent code" ),
				clEnumValN( llvm::Reloc::DynamicNoPIC, "dynamic-no-pic",
							"Relocatable external references, non-relocatable code" ),
				clEnumValN( llvm::Reloc::ROPI, "ropi", "Code and read-only data relocatable, accessed PC-relative" ),
				clEnumValN( llvm::Reloc::RWPI, "rwpi",
							"Read-write data relocatable, accessed relative to static base" ),
				clEnumValN( llvm::Reloc::ROPI_RWPI, "ropi-rwpi", "Combination of ropi and rwpi" ) ),
	cl::cat( category ) );
cl::opt<CodeGenFileType> FileType( "filetype",
								   cl::desc( "Choose a file type (not all types are supported by all targets):" ),
								   cl::init( CGFT_ObjectFile ),
								   cl::values( clEnumValN( CGFT_AssemblyFile, "asm", "Emit an assembly ('.s') file" ),
											   clEnumValN( CGFT_ObjectFile, "obj", "Emit a native object ('.o') file" ),
											   clEnumValN( CGFT_Null, "null", "Emit nothing, for performance testing" ),
											   clEnumValN( CGFT_IR, "ir", "Emit LLVM IR ('.ll') file" ) ),
								   cl::cat( category ) );

static llvm::TargetOptions InitTargetOptionsFromCodeGenFlags() {
	llvm::TargetOptions Options;

	Options.DataSections = DataSections;

	// Options.MCOptions = InitMCTargetOptionsFromFlags();

	return Options;
}

static llvm::StringRef getOutputFilename( char const* targetName, llvm::Triple const& triple ) {
	if ( !OutputFilename.empty() ) {
		return OutputFilename;
	}

	if ( InputFilename == "-" ) {
		OutputFilename = "-";
		return OutputFilename;
	}

	OutputFilename = llvm::StringRef( InputFilename );
	switch ( FileType ) {
	default:
	case CGFT_AssemblyFile:
		assert( targetName );
		if ( targetName[0] == 'c' ) {
			if ( targetName[1] == 0 ) {
				OutputFilename += ".cbe.c";
			} else if ( targetName[1] == 'p' && targetName[2] == 'p' ) {
				OutputFilename += ".cpp";
			} else {
				OutputFilename += ".s";
			}
		} else {
			OutputFilename += ".s";
		}
		break;

	case CGFT_ObjectFile:
		if ( triple.getOS() == llvm::Triple::Win32 ) {
			OutputFilename += ".obj";
		} else {
			OutputFilename += ".o";
		}
		break;

	case CGFT_Null:
		OutputFilename += ".null";
		break;

	case CGFT_IR:
		OutputFilename += ".ll";
		break;
	}
	return OutputFilename;
}

static llvm::StringRef getVariableName() {
	if ( !VariableName.empty() ) {
		return VariableName;
	}

	if ( InputFilename == "-" ) {
		fprintf( stderr, "warn: variable name not specified, and could not be guessed\n" );
		return "stdin";
	}

	VariableName = llvm::sys::path::filename( InputFilename );
	std::replace( VariableName.begin(), VariableName.end(), '.', '_' );
	return VariableName;
}

static std::string getTargetTriple() {
	if ( !TargetTriple.empty() ) {
		return llvm::Triple::normalize( TargetTriple );
	}

	return llvm::sys::getDefaultTargetTriple();
}

static llvm::StringRef getCPUStr() {
	// If user asked for the 'native' CPU, autodetect here. If autodection fails,
	// this will set the CPU to an empty string which tells the target to
	// pick a basic default.
	if ( MCPU == "native" ) {
		return llvm::sys::getHostCPUName();
	}

	return MCPU.getValue();
}

static std::string getFeaturesStr() {
	llvm::SubtargetFeatures Features;

	// If user asked for the 'native' CPU, we need to autodetect features.
	// This is necessary for x86 where the CPU might not support all the
	// features the autodetected CPU name lists in the target. For example,
	// not all Sandybridge processors support AVX.
	if ( MCPU == "native" ) {
		llvm::StringMap<bool> HostFeatures;
		if ( llvm::sys::getHostCPUFeatures( HostFeatures ) ) {
			for ( auto& F : HostFeatures ) {
				Features.AddFeature( F.first(), F.second );
			}
		}
	}

	for ( unsigned i = 0; i != MAttrs.size(); ++i ) {
		Features.AddFeature( MAttrs[i] );
	}

	return Features.getString();
}

#if LLVM_VERSION >= 600
static llvm::Optional<llvm::CodeModel::Model> getCodeModel() {
	if ( CMModel.getNumOccurrences() ) {
		llvm::CodeModel::Model M = CMModel;
		return M;
	}
	return llvm::None;
}
#else
static llvm::CodeModel::Model getCodeModel() {
	if ( CMModel.getNumOccurrences() ) {
		return CMModel;
	}
	return llvm::CodeModel::Default;
}
#endif

static llvm::Optional<llvm::Reloc::Model> getRelocModel() {
	if ( RelocModel.getNumOccurrences() ) {
		llvm::Reloc::Model R = RelocModel;
		return R;
	}
	return llvm::None;
}

static llvm::TargetMachine::CodeGenFileType asTargetMachineCodeGenFileType( CodeGenFileType filetype ) {
	assert( filetype != CGFT_IR );
	return static_cast<llvm::TargetMachine::CodeGenFileType>( filetype );
}

int main( int argc, char const* argv[] ) {
// Initialize LLVM.
#if LLVM_CONFIG >= 700
	llvm::InitLLVM X( argc, argv );
#endif

	// Initialize targets first, so that --version shows registered targets.
	llvm::InitializeAllTargets();
	llvm::InitializeAllTargetMCs();
	llvm::InitializeAllAsmPrinters();
	llvm::InitializeAllAsmParsers();

	// Process the command line.
	cl::HideUnrelatedOptions( category );
	if ( !cl::ParseCommandLineOptions( argc, argv ) ) {
		return EXIT_FAILURE;
	}

	// Initalize the module.
	llvm::LLVMContext context;
	llvm::Module module( "binc", context );
	module.setSourceFileName( InputFilename );
	module.setTargetTriple( getTargetTriple() );
	auto const bufferSize = buildModule( module, getVariableName(), NullTerminate );
	assert( !module.global_empty() );

	// Initialize debug information
	if ( DebugInfo ) {
		buildDebugInfo( module, bufferSize );
	}

	if ( FileType == CGFT_IR ) {
		std::error_code ec;
		auto triple = llvm::Triple( module.getTargetTriple() );
		auto filename = getOutputFilename( "", triple );
		llvm::raw_fd_ostream out( filename, ec, llvm::sys::fs::OpenFlags::F_None );
		if ( ec ) {
			fprintf( stderr, "error: could not open output: %s\n", ec.message().c_str() );
			return EXIT_FAILURE;
		}
		module.print( out, nullptr );
		return EXIT_SUCCESS;
	}

	std::string errMessage;
	auto triple = llvm::Triple( module.getTargetTriple() );
	auto const* target = llvm::TargetRegistry::lookupTarget( MArch, triple, errMessage );
	if ( !target ) {
		fprintf( stderr, "error: could initialize the target: %s\n", errMessage.c_str() );
		return EXIT_FAILURE;
	}

	auto options = InitTargetOptionsFromCodeGenFlags();
	std::unique_ptr<llvm::TargetMachine> targetMachine(
		target->createTargetMachine( triple.getTriple(), getCPUStr(), getFeaturesStr(), options, getRelocModel(),
									 getCodeModel(), llvm::CodeGenOpt::None ) );
	if ( !targetMachine ) {
		fprintf( stderr, "error: could not allocate target machine\n" );
		return EXIT_FAILURE;
	}

	auto ec = writeOutputFile( module, getOutputFilename( target->getName(), triple ),
							   asTargetMachineCodeGenFileType( FileType ), *targetMachine );
	if ( ec ) {
		fprintf( stderr, "error: could open output: %s\n", ec.message().c_str() );
		return EXIT_FAILURE;
	}

	return EXIT_SUCCESS;
}
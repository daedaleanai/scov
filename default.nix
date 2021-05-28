{ buildGoModule ? (import <nixpkgs> {}).buildGoModule,
  lib ? (import <nixpkgs> {}).lib }:

buildGoModule rec {
  pname = "scov";
  version = "v0.9.1";

  src = ./.;

  # No need to build subpackage behemoth, which is only for testing.
  subPackages = [ "." ];

  vendorSha256 = "18miyyil4jpmf3v1axkn3k1lhza07p9p26agvmqi8mlwkraabhxb";

  # Update the version information in the built executable
  buildFlagsArray = ''
     -ldflags=
         -X main.versionInformation=${version}
         -s -w
  '';

  meta = with lib; {
    description = "Generate reports on code coverage using gcov, lcov, or llvm-cov.";
    homepage = "https://gitlab.com/stone.code/scov";
    license = licenses.bsd3;
    platforms = platforms.all;
  };
}

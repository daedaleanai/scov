{ buildGoModule ? (import <nixpkgs> {}).buildGoModule,
  lib ? (import <nixpkgs> {}).lib }:

buildGoModule rec {
  pname = "scov";
  version = "v0.9.1";

  src = ./.;

  # No need to build subpackage behemoth, which is only for testing.
  subPackages = [ "." ];

  vendorSha256 = "0c28d6dip04m0hljss9llp84nc1a0l7vc11zyxlidnpgd5kychxp";

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

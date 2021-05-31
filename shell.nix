{ pkgs ? import <nixpkgs> { }, enableLint ? false }:

pkgs.mkShell {
  # No dependencies beyond stdlib.
  buildInputs = [ pkgs.go pkgs.gitMinimal ]
    ++ (if enableLint then [ pkgs.golangci-lint ] else [ ]);

  # Make sure we don't pick up the users' GOPATH.
  # Advertise the current version of go when shell starts.
  shellHook = "unset GOPATH; go version;";
}

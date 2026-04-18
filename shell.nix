with import <nixpkgs> {};
  stdenv.mkDerivation {
    name = "go";
    buildInputs = [
      go
      gotestsum
      gofumpt
      air
    ];
    shellHook = ''
      export PATH="$HOME/go/bin:$PATH"
    '';
  }

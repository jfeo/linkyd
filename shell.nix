{pkgs ? import <nixpkgs> {}}:
pkgs.mkShell {
  name = "go";
  packages = with pkgs; [
    go
    gotestsum
    gofumpt
    air
  ];
  shellHook = ''
    export PATH="$HOME/go/bin:$PATH"
  '';
}

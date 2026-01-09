{
  description = "🤖 Telegram notifications for Claude Code hooks";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages.default = pkgs.buildGoModule {
          pname = "tgnotify";
          version = "0.1.0";
          src = ./.;
          vendorHash = "sha256-pbA/AlBz3cQYRTMnQ/qBPcinYOKokrBLNhkbRTq54gE=";
          meta = with pkgs.lib; {
            description = "Telegram notifications for Claude Code hooks";
            homepage = "https://github.com/nexo-tech/tgnotify";
            license = licenses.mit;
            mainProgram = "tgnotify";
          };
        };

        devShells.default = pkgs.mkShell {
          buildInputs = [ pkgs.go ];
        };
      }
    );
}

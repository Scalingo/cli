{
  inputs = {
    nixpkgs.url = "nixpkgs";
    utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, utils }:
    utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in rec {
        packages = rec {
          default = pkgs.buildGoModule {
            name = "scalingo";

            src = pkgs.fetchFromGitHub {
              owner = "Scalingo";
              repo = "cli";
              rev = "1.28.2";
              sha256 = "sha256-dMiOGPQ2wodVdB43Sk3GfEFYIU/W2K9DG/4hhVxb1fs=";
            };

            vendorSha256 = null;

            ldflags = [ "-w" "-s" ];

            preConfigure = ''
              export HOME=$TMPDIR
            '';
          };
        };

        apps = rec {
          default = utils.lib.mkApp { drv = packages.default; };
        };

        devShell = with pkgs;
          mkShell {
            buildInputs = [ go ];
          };
      });
}

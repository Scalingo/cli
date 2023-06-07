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
              rev = "1.29.1";
              sha256 = "sha256-xBf+LIwlpauJd/0xJIQdfEa0rxph3BJPuMY4+0s+Bb4=";
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

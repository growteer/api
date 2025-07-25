{
  description = "Growteer Backend Flake";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  inputs.flake-utils.url = "github:numtide/flake-utils";
  inputs.gomod2nix.url = "github:nix-community/gomod2nix";
  inputs.gomod2nix.inputs.nixpkgs.follows = "nixpkgs";
  inputs.gomod2nix.inputs.flake-utils.follows = "flake-utils";

  outputs =
    { self
    , nixpkgs
    , flake-utils
    , gomod2nix
    ,
    }:
    (flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs {
          config.allowUnfree = true;
          inherit system;
        };

        # The current default sdk for macOS fails to compile go projects, so we use a newer one for now.
        # This has no effect on other platforms.
        callPackage = pkgs.darwin.apple_sdk_11_0.callPackage or pkgs.callPackage;

      in
      rec {
        packages.default = callPackage ./. {
          inherit (gomod2nix.legacyPackages.${system}) buildGoApplication;
        };
        packages.container = callPackage ./container.nix {
          package = packages.default;
        };
        devShells.default = callPackage ./shell.nix {
          inherit (gomod2nix.legacyPackages.${system}) mkGoEnv gomod2nix;
        };
      }
    ));
}

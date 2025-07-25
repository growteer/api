{ pkgs ? (
    let
      inherit (builtins) fetchTree fromJSON readFile;
      inherit ((fromJSON (readFile ./flake.lock)).nodes) nixpkgs gomod2nix;
    in
    import (fetchTree nixpkgs.locked) {
      overlays = [
        (import "${fetchTree gomod2nix.locked}/overlay.nix")
      ];
    }
  )
, buildGoApplication ? pkgs.buildGoApplication
}:

buildGoApplication {
  pname = "growteer-api";
  version = "0.1";
  src = ./.;
  subPackages = [ "cmd/growteer-api" ];
  modules = ./gomod2nix.toml;

  # NOTE (c.nicola): disabled, run `nix develop -c make test` instead
  doCheck = false;

  CGO_ENABLED = 0;

  ldflags = [ "-s -w" ];
}

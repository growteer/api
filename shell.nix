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
, mkGoEnv ? pkgs.mkGoEnv
, gomod2nix ? pkgs.gomod2nix
, mongodb ? pkgs.mongodb
, mongosh ? pkgs.mongosh
, gqlgen ? pkgs.gqlgen
, golangci-lint ? pkgs.golangci-lint
, go-mockery ? pkgs.go-mockery
, go-junit-report ? pkgs.go-junit-report
, gopls ? pkgs.gopls
}:

let
  goEnv = mkGoEnv { pwd = ./.; };
in
pkgs.mkShell {
  packages = [
    goEnv
    gomod2nix
    mongodb
    mongosh
    gqlgen
    golangci-lint
    go-mockery
    go-junit-report
    gopls
  ];

  env = {
    MONGODB_URI = "mongodb://localhost:27017/?directConnection=true";
    MONGODB_DATABASE = "growteer";
    JWT_SECRET = "baconkilbasa";
    ALLOWED_ORIGINS = "http://localhost:5173";
    HTTP_PORT = "8080";
    SESSION_TTL_MINUTES = "60";
  };

  shellHook = ''
    # Create necessary directories
    mkdir -p ./data/db
    mkdir -p ./data/log

    # Start MongoDB as a service
    echo "Starting MongoDB..."
    mongod --quiet --fork --logpath ./data/log/mongod.log --dbpath ./data/db

    # Cleanup command to stop MongoDB when exiting the shell
    trap 'echo "Stopping MongoDB..."; mongod --quiet --shutdown --dbpath ./data/db' EXIT

    echo "Connect to mongodb with:"
    echo "Connection URI: $MONGODB_URI"
    echo "Database Name: $MONGODB_DATABASE"
    echo ""
    echo "JWS secret: $JWT_SECRET"
  '';
}

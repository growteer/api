[![codecov](https://codecov.io/gh/growteer/api/graph/badge.svg?token=Y6SAD10060)](https://codecov.io/gh/growteer/api)

# Setup

- Install Nix, instructions are found here [nixos.org](https://nixos.org/download/)
- Launch a development environment with `nix develop`
- Build the project with `nix build` (output path is `result/bin`)
- Test the project with `nix develop -c make test` (or run `make test` while inside the development shell)
- Lint the project with `nix develop -c make lint` (or run `make lint` while inside the development shell)
- Build a docker image with `nix build .#container` and load the image with `docker load < result`

All current dependencies are ready and available inside `nix develop`, so it is encouraged to use it for developing.

# Upgrade

If the `go.mod` file changes due to an upgrade, please remember to run `gomod2nix generate` to rehydrate the `gomod2nix.toml` file. This file
tracks the dependencies as separatate nix packages and are thus versioned with the flake.

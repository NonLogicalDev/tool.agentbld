{
  description = "Render and sync layered agent configuration";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";
    flake-parts.inputs.nixpkgs-lib.follows = "nixpkgs";
  };

  outputs = inputs@{ self, nixpkgs, flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = nixpkgs.lib.systems.flakeExposed;

      perSystem = { config, pkgs, ... }: {
        packages.agentbld = pkgs.callPackage ./nix/package.nix {
          revision = self.rev or (self.dirtyRev or null);
        };
        packages.default = config.packages.agentbld;
      };
    };
}

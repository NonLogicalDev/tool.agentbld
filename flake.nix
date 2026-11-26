{
  description = "Render and sync layered agent configuration";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";

  outputs =
    { self, nixpkgs }:
    let
      inherit (nixpkgs) lib;
    in
    {
      packages = lib.genAttrs lib.systems.flakeExposed (
        system:
        let
          agentbld = nixpkgs.legacyPackages.${system}.callPackage ./nix/package.nix {
            revision = self.rev or (self.dirtyRev or null);
          };
        in
        {
          inherit agentbld;
          default = agentbld;
        }
      );
    };
}

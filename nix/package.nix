{
  lib,
  buildGoModule,
  source ? lib.fileset.toSource {
    root = ../.;
    fileset = lib.fileset.unions [
      ../cmd
      ../internal
      ../go.mod
      ../go.sum
    ];
  },
  version ? "0.0.0-unstable",
  revision ? null,
}:

buildGoModule {
  pname = "agentbld";
  inherit version;
  src = source;
  vendorHash = "sha256-vkXeJTYk088+AfRF0LV9HJIdWUt6By8VtU1Ve5D33ls=";
  subPackages = [ "cmd/agentbld" ];

  env.CGO_ENABLED = "0";
  ldflags = [
    "-s"
    "-w"
    "-X main.version=${version}"
  ] ++ lib.optional (revision != null) "-X main.commit=${revision}";

  checkPhase = ''
    runHook preCheck
    go test ./...
    runHook postCheck
  '';

  meta = {
    description = "Render and sync layered agent configuration";
    homepage = "https://github.com/NonLogicalDev/tool.agentbld";
    mainProgram = "agentbld";
    platforms = lib.platforms.darwin ++ lib.platforms.linux ++ lib.platforms.freebsd;
  };
}

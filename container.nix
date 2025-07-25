{ pkgs, package }:


pkgs.dockerTools.buildImage {
  name = "growteer-api";
  tag = "0.1";
  created = "now";
  copyToRoot = pkgs.buildEnv {
    name = "image-root";
    paths = [ package ];
    pathsToLink = [ "/bin" ];
  };
  config.Cmd = [ "${package}/bin/growteer-api" ];
}

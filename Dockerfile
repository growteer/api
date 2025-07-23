FROM nixos/nix:latest AS builder

COPY . /tmp/build

WORKDIR /tmp/build

RUN nix-env -iA nixpkgs.cacert

RUN cp /etc/ssl/certs/ca-bundle.crt ca-certificates
RUN nix --extra-experimental-features "nix-command flakes" build
RUN mkdir /tmp/nix-store-closure
RUN cp -R $(nix-store -qR result/) /tmp/nix-store-closure

FROM scratch

WORKDIR /bin

COPY --from=builder /tmp/nix-store-closure /nix/store
COPY --from=builder /tmp/build/ca-certificates /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /tmp/build/result/bin/growteer-api .

CMD ["/bin/growteer-api"]

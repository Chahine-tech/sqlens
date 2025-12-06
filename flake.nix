{
  description = "SQLens - High-performance multi-dialect SQL query analysis tool";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};

        # Package version
        version = "0.1.0";

      in
      {
        # Development shell with all tools
        devShells.default = pkgs.mkShell {
          name = "sql-parser-go-dev";

          buildInputs = with pkgs; [
            # Go toolchain (latest stable)
            go

            # Go development tools
            gopls              # Language server
            gotools            # goimports, godoc, etc.
            go-tools           # staticcheck, etc.
            golangci-lint      # Linter
            delve              # Debugger

            # Build tools
            gnumake

            # Optional: useful for benchmarking
            hyperfine
          ];

          shellHook = ''
            echo "🚀 SQL Parser Go - Development Environment"
            echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
            echo "Go version: $(go version)"
            echo ""
            echo "Available commands:"
            echo "  make build    - Build the project"
            echo "  make test     - Run tests"
            echo "  make bench    - Run benchmarks"
            echo "  make lint     - Run linter"
            echo "  make fmt      - Format code"
            echo ""
            echo "Run 'make help' for more commands"
            echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
          '';

          # Environment variables
          CGO_ENABLED = "0";
        };

        # Main package - sqlparser binary
        packages.default = pkgs.buildGoModule {
          pname = "sqlparser";
          inherit version;

          src = ./.;

          # Vendor hash for Go dependencies
          # Computed automatically by Nix
          vendorHash = "sha256-g+yaVIx4jxpAQ/+WrGKxhVeliYx7nLQe/zsGpxV4Fn4=";

          # Build flags
          ldflags = [
            "-s"
            "-w"
            "-X main.version=${version}"
          ];

          # Run tests
          doCheck = true;
          checkPhase = ''
            runHook preCheck
            go test -v ./...
            runHook postCheck
          '';

          meta = with pkgs.lib; {
            description = "High-performance SQL parser with multi-dialect support";
            homepage = "https://github.com/Chahine-tech/sql-parser-go";
            license = licenses.mit;
            maintainers = [ ];
            mainProgram = "sqlparser";
          };
        };

        # Alias for the package
        packages.sqlparser = self.packages.${system}.default;

        # Apps - make it easy to run
        apps.default = {
          type = "app";
          program = "${self.packages.${system}.default}/bin/sqlparser";
        };

        # Checks - run tests with nix flake check
        checks = {
          # Run all tests
          tests = pkgs.runCommand "sql-parser-tests" {
            buildInputs = [ pkgs.go ];
          } ''
            cp -r ${./.}/* .
            chmod -R +w .
            export HOME=$TMPDIR
            go test -v ./tests
            touch $out
          '';

          # Check formatting
          fmt-check = pkgs.runCommand "sql-parser-fmt-check" {
            buildInputs = [ pkgs.go ];
          } ''
            cp -r ${./.}/* .
            chmod -R +w .

            # Check if code is formatted
            unformatted=$(gofmt -l .)
            if [ -n "$unformatted" ]; then
              echo "The following files are not formatted:"
              echo "$unformatted"
              exit 1
            fi

            touch $out
          '';
        };

        # Formatter for `nix fmt`
        formatter = pkgs.nixpkgs-fmt;
      }
    );
}

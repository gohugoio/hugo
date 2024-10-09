#!/bin/sh

# Check if a custom hugo-docker-entrypoint.sh file exists.
if [ -f hugo-docker-entrypoint.sh ]; then
  # Execute the custom entrypoint file.
  sh hugo-docker-entrypoint.sh "$@"
  exit $?
fi

# Check if a package.json file exists.
if [ -f package.json ]; then
  # Check if node_modules exists.
  if [ ! -d node_modules ]; then
    # Install npm packages.
    # Note that we deliberately do not use `npm ci` here, as it would fail if the package-lock.json file is not up-to-date,
    # which would be the case if you run the container with a different OS or architecture than the one used to create the package-lock.json file.
    npm i
  fi
fi

exec "hugo" "$@"
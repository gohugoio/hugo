docker run --rm -v .:/project -v $HOME/Library/Caches/hugo_cache:/cache  -u $(id -u):$(id -g) testing_bullseye build


docker build -t testing_alpine .

                   
docker run --rm -v .:/project -v $HOME/Library/Caches/hugo_cache:/cache -p 1313:1313 testing_alpine server --bind="0.0.0.0"


This is the entry point file:

```sh
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
    npm i
  fi
fi

exec "hugo" "$@"
```

If a custom `hugo-docker-entrypoint.sh` script exists in the root of the Hugo project, that script will be executed instead of the default entrypoint script. 
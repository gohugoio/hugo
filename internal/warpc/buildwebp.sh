pushd genwebp
make || exit 1
popd
touch webp_integration_test.go
gobench -package=./hugolib -bench="BenchmarkSiteBuilding/YAML,num_langs=3,num_pages=5000,tags_per_page=5,shortcodes,render" -count=3 > 1.bench
benchcmp -best 0.bench 1.bench
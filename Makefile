.PHONY: test test-js test-bash

test: test-js test-bash

test-js:
	npx vitest run

test-bash:
	npx bats src/__tests__/*.bats

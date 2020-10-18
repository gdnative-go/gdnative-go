# CHANGELOG

<!--- next entry here -->

## 0.1.0
2020-10-18

### Features

- removed Python generation scriopt and move build to Magefile (7c2d111571a3ab931b6e2f22e0d21c1c66e34366)
- add capabilities for code generation with autoregistration (7cc63f67d4771a2792d3de1750293c611d23e9b3)

### Fixes

- added more missing comments (8498ef8aa28a4b8018da5052485e08c78a505030)
- fallback to goreturns if goimport is missing (56c65385ad09478ce932d4f82b4d1895e04293db)
- fix wrong conversion to uintptr to apply pointer aritmetict to **C.godot_variant to iterate the array (a00ff223d5832533c5d0e9d4847fbb1db13d04f3)
- fix magefile clean (30e084b98292d067da5159c8208024949d92481f)
- fixed errors and added more examples (552d344e17682f79f6e2f8602500f6010c0c5d09)
- add clearer docstrings to SimpleDemo (0912bec369556d5f3cee7abcaca745cb0bb91e8a)
- add clearer docstrings on SimpleDemoManual (94b61cc463e3783d51f2f7bf6fd6344c7b7a55c5)
- updated README.md (10aa39fc5d7bf708820ae4fd48abe25e000b656a)
- removed invalid tag (33551ec115b80f7e3b286bdf0c9e3fb37f070c08)
- fix compilation step on test stage (4a20358164b17412ac73b1b5ac1d5fd823aed973)
- initialize git submodules before trying to build the library (b5311d4182da495e1e61992b36e0db4c6e248a3d)
- install goreturns needed for generation (de4a478a483494ff5346606ae7694770b4b19837)
- fix issues all around gdnative (9a82e5c4250f804d6776a24fe6567ebc01886e4b)
- eerything working in auto generation but signals make editor to crash (probably due pointer conversion into array) (995747d331ae10aaecf887ea8590baed8bc7bc7e)
- fixed signals and removed unused imports (0f2165e175a634284e20b91bdbb3591da87c2fca)
- fix linting errors and rename gdnativego compiler into gogdc (8419ba6c1c9ce8e66c65524f6a18d1a4e28753c5)
- fix generation bugs introduced by sloppy linting errors refactor (7336dbcb9c093157bead0f1d6e4bb4d6c63b71dd)
- fix generate pipeline stage (82e61b1c2d6e983c6489934abb35e07489483ecf)
- simplified pipeline steps as some are dependent on others (10fdc146675e0dcb675278de46a5a796c25e1d62)
- change goreturns invocation due Go caching (8051fbd3ede6d0e5d2d618d4873a7a94e77d252d)
- fix README.md typos (7fffdb0ad926b179e6f0bd24e024c7d683a1ed99)
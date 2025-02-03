# Changelog
All notable changes are recorded here.

### Format

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/).

Entries should have the imperative form, just like commit messages. Start each entry with words like
add, fix, increase, force etc.. Not added, fixed, increased, forced etc.

Line wrap the file at 100 chars.                                              That is over here -> |

### Categories each change fall into

* **Added**: for new features.
* **Changed**: for changes in existing functionality.
* **Deprecated**: for soon-to-be removed features.
* **Removed**: for now removed features.
* **Fixed**: for any bug fixes.
* **Security**: in case of vulnerabilities.


## [1.1.2] - 2025-02-03
### Changed
- Use Go 1.23.5
### Security
- Update third-party dependencies.


## [1.1.1] - 2024-09-06
### Fixed
- Do not share the XOR key offset between the send and receive threads.


## [1.1.0] - 2024-09-05
### Added
- Add XOR v2.


## [1.0.4] - 2024-07-04
### Changed
- Upgrade to use Go 1.22.5
- Cancel forwarding swiftly if either side terminates the connection.
### Security
- Update vulnerable dependencies from stdlib.


## [1.0.3] - 2024-06-18
### Changed
- Upgrade to use Go v1.22.4
- Create zip-artifacts deterministically.
### Security
- Update vulnerable dependencies from stdlib.
- Update third-party dependencies.


## [1.0.2] - 2024-06-03
### Added
- Add build target for Windows ARM64.


## [1.0.1] - 2024-04-19
### Added
- Add build target for Linux ARM64.
### Changed
- Downgrade go v1.22.2 -> v1.21.3.


## [1.0.0] - 2024-04-18
### Added
- Core functionality with some configurability.

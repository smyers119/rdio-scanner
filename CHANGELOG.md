# Change log

## Version 6.4

- New `-cmd` command line options to allow advanced administrative tasks.
- New playback mode goes live options which is not enabled by default (issue #175).
- Fix logs retrieval from administrative dashboard (issue #193).
- Improved field conversions when retrieving calls from a mysql/mariadb database (issue #194, #198).
- Highlight replayed call on the history list (issue #196).

## Version 6.3

- Change scroll speed when drag droping talkgroups or units in a system (discussion #170).
- System Ids listed in the `Config / Options / AFS Systems` will have their talkgroup Ids displayed in AFS format (issue #163).
- New dirwatch meta tags #GROUP #SYSLBL #TAG #TGAFS and #UNIT for better ProScan compatibility (issue #164).
- Playback mode will now catch up to live (issue #175).
- Dirwatch code rewrite (issue #177).

_v6.3.1_

- Playback mode catch up to live, then swith to livefeed mode.
- Removed the mutex lock on Clients.Count which led to a deadlock and froze call ingestion.

_v6.3.2_

- New #TGLBL metatag for dirwatch for ProScan (%C) or alike.
- Fix `semacquire` lockup in Clients (issue #177, #181, #182).
- Replay button now replays from history if pressed multiple times quickly (issue #186).

_v6.3.3_

- Fix concurrent map writes fatal error in dirwatch (issue #187).
- Brighter LED colors and new orange color.
- Fix call id when retrieved from a MySQL database.
- Add loudnorm audio filter to the ffmpeg audio conversion.
- Show the real IP address in the logs taking into account if behind a proxy.
- Fix panic when emitting a call to clients.

_v6.3.4_

- Fix ffmpeg audio filter not available on older version (issue #189).
- Improved logging when run as a service, Windows users can now see these logs in the events viewer.
- Dirwatch now catches panic errors and logs them.

_v6.3.5_

- Replace standard map with sync.map in dirwatch.
- Fix the ffmpeg version test.
- Fix led color type, orage -> orange.
- Fix incorrect options when reading from a mysql database (issue #190).

_v6.3.6_

- Fix systems order properties not sent to clients.
- Fix side panels not scrolling to top when opened.

## Version 6.2

- New max clients options which is 200 by default.
- New show listeners count options which is disabled by default (issue #125).
- Fix panic: concurrent write to websocket connection on goroutine.
- Fix units import from SDR Trunk (issue #150).

_v6.2.1_

- Fix SIGSEGV error in Units.Merge (issue #151).

_v6.2.2_

- Fix another SIGSEGV error in Units.Merge (issue #151).

_v6.2.3_

- New random UUID in the JSON-Web Token payload.
- Fix dirwatch not properly shutting down when a new configuration is applied.
- Fix dashboard logout not sending HTTP 200 OK status.
- Clear the active dirwatch list when stopped.
- Pauses calls ingestion before database pruning.
- Fix regex for units in driwatch type SDRTrunk (discussion #155).
- Update SQLite driver.

_v6.2.4_

- Fix call frequencies table not being transmitted to downstream.
- Avoid using setInterval and setTimeout in the webapp.
- Fix talkgroup search filter upon new configuration (issue #158).

_v6.2.5_

- Fix unnecessary auto populate of unit id/label (issue #160).

## Version 6.1

- Calls now support patched talkgroups.
- New search patched talkgroups option which is disabled by default.
- Talkgroups and units are now stored in their own database table.
- New units CSV importer.
- Fix blacklisted talkgroups being created anyway when autopopulate is enabled.
- Fix compatibility with mysql/mariadb (default sqlite is still recommended).

_v6.1.1_

- Fix `unknown datetime format sql.NullString` error.

_v6.1.2_

- Fix image links in webapp.md (issue #76).
- Fix SIGSEGV when trying to autopopulate (issue #77).
- Fix parsing SDRTrunk meta data.
- Dirwatch type trunk-recorder now deletes json files without audio (when deleteAfter is set).
- Add a new `docs/update-from-v5.md` document.

_v6.1.3_

- Fix concurrent config write when autopopulate is enabled (issue #77).
- Fix API in regards to audio filename and audio type (issue #78).
- Fix migration error on mysql database (issue #86).
- Fix some calls not playing on the native app (issue #87).
- Fix admin password not read from mysql.

_v6.1.4_

- Talkgroup label now syncs with the talkgroup_tag from the API or dirwatch (issue #80).
- Fix more migration errors on mysql database (issue #86).
- Fix config export not working with non latin-1 characters (issue #89).
- Fix talkgroup label from dirwatch type sdrtrunk (discussion #98).
- Fix SIGSEGV (issue #100).
- New `patch` indicator for patched talkgroups.

_v6.1.5_

- Fix trunk-recorder API (issue #104).
- Fix for avoid/patch flags on main display not beaving as expected.
- Fix downstream not sending sources data.
- Fix dirwatch crashing when config is updated.

_v6.1.6_

- Fix webapp not reporting the correct version.

_v6.1.7_

- More concurrency mutexes to resolve SQL_BUSY errors.
- Better internal management of dirwatches.
- Fix SDRTrunk files not being ingested (discussion #108).
- Fix Trunk Recorder talkgroup_tag assign to the wrong property (issue #115).
- Improved the way the talkgroup label and name are autopopulated. If Trunk Recorder sends a talkgroup_tag with an empty value or with a single `-`, it will not overwrite the talkgroup label.

_v6.1.8_

- New dirwatch masks #TGHZ, #TGKHZ and #TGMHZ which allow to set talkgroup id based on frequency.

_v6.1.9_

- Fix talkgroup sorting issue when importing from a CSV file (issue #119).
- Fix SIGSEGV (issue #120).

_v6.1.10_

- Backport dirwatch delay value from v5.1.

_v6.1.11_

- Fix connection errors when behind a reverse-proxy.
- Fix disappearing talkgroups (issue #127).

_v6.1.12_

- Fix too many open files (issue #129).
- Cosmetic: AVOID and PATCH flags now only appear when needed.

_v6.1.13_

- Better handling of dead client connections.
- Fix too many open files (issue #129).
- Remove net.http error messages from the output (issue #131).

_v6.1.14_

- Fix FAQ section not being added to the PDF documents.
- Bump delay before killing unauthenticated clients from 10 seconds to 60 seconds.
- Remove the gitter.im support forum from the documentation and prefer github discussions.

_v6.1.15_

- Fix access and downstreams order not retained.
- Remove the self-signed certificate generator (-ssl create) as it was causing more problems than solutions.
- Client handling and call ingestion now run on 2 different threads (issue #135).
- Fix downstream talkgroup select keeps reverting to all talkgroups (issue #136).

_v6.1.16_

- Fix concurrent map access for clients.
- Some tweaks to websocket management.

## Version 6.0

- Backend server rewritten in Go for better performance and ease of installation.
- New toggle by tags option to toggle talkgroups by their tag in addition to their group.
- Buttons on the select panel now sound differently depending on their state.
- You can now filter calls by date and time on the search panel.
- Installable as a service from the command line.
- Let's Encrypt automatic generation of certificates from the command line.
- A bunch of minor fixes and improvements.

### BREAKING CHANGES SINCE V5

[Rdio Scanner](https://github.com/chuot/rdio-scanner) is now distributed as a precompiled executable in a zip file, which also contains documentation on how it works.

The backend server has been completely rewritten in GO language. Therefore, all the subpackages used in v5 had to be replaced with new ones. These new subpackages do not necessarily have the same functionality as those of v5.

- No more polling mode for _dirwatch_, which in a way is a good thing as polling was disastrous for CPU consumption. The alternative is to install a local instance and use the downstream feature to feed your main instance.
- Due to the polling situation, the Docker version of Rdio Scanner doesn't have the dirwatch feature.
- Default database name changed from _database.sqlite_ to _rdio-scanner.db_. You will need to rename your database file with the new name if you want to convert it. Otherwise, a new database will be created.

_v6.0.1_

- Fix button sound on select panel for TG (beep state inverted)
- Auto populate system units (issue #66)

_v6.0.2_

- Try to fix the SQL_BUSY error (issue #67).
- Fix `-service stop` timing out before exiting.
- Drop the ApiKey uniqueness of the downstreams database table.
- Fix auto-populating the database with empty units tag.

_v6.0.3_

- Fix strconv.Atoi: invalid syntax for dirwatch type sdrtrunk.
- Fix the new version available dialog opening more than once.

_v6.0.4_

- Fix wrong time calculation in prune scheduler.
- More fix on the SQL_BUSY error (issue #67).
- Support files (certs, db, ini) are now created in the same folder as the executable, if the folder is writable, or under a `Rdio Scanner` folder in the user's home folder.
- Some code refactoring.

_v6.0.5_

- Force mime type to `application/javascript` for `.js` extensions. (see https://github.com/golang/go/issues/32350).
- New `-base_dir` option to specify the directory where all data will be written.
- New Docker container with disabled dirwatch.

_v6.0.6_

- Fix an issue with not closing the database when restarting the host platform (issue #71).
- Fixed SDRTunk parser when artist tag contains CTCSS tones.
- Platforms linux/amd64, linux/arm and linux/arm64 are now available for the Docker container.

_v6.0.7_

- Fix dropped connections when going through a proxy.

## Version 5.2

- Change to how the server reports version.
- Fix cmd.js exiting on inexistant session token keystore.
- Fix issue with iframe.
- Node modules updated for security fixes.

_v5.2.1_

- Fix talkgroup header on the search panel (issue #47).
- Update dirwatch meta tags #DATE, #TIME and #ZTIME for SDRSharp compatibility (issue #48).
- Fix dirwath date and time parsing bug.
- Configurable call duplicate detection time frame.

_v5.2.2_

- Little changes to the main screen history layout, more room for the second and third columns.
- Node modules updates.

_v5.2.3_

- Change history columns padding from 1px to 6px on the main screen.
- Fix a bug in the admin api where the server crash when saving new config from the admin dashboard.

_v5.2.4_

- Updated to Angular 12.2.
- New update prompt for clients when server is updated.
- Fix unaligned back arrow on the search panel.

_v5.2.5_

- STS command removed from the server.
- Minor fixes here and there.
- README.md updated.
- Documentation images resized.

_v5.2.6_

- Fix crash when when options.pruneDays = 0.

_v5.2.7_

- Fix handling of JSON datatypes on MySQL/MariaDB database backend.
- Fix listeners count.

_V5.2.8_

- Fix SQLite does not support TEXT with options.

_V5.2.9_

- Fix bad code for server options parsing.
- Increase dirwatch polling interval from 1000ms to 2500ms.

## Version 5.1

This one is a big one... **Be sure to backup your config.json and your database.sqlite before updating.**

- With the exception of some parameters like the SSL certificates, all configurations have been moved to an administrative dashboard for easier configuration. No more config.json editing!
- Access codes can now be set with a limit of simultaneous connections. It is also possible to configure an expiration date for each access codes.
- Auto populate option can now be set per system in addition to globally.
- Duplicate call detection is now optional and can be disabled from the options section of the administrative dashboard.
- On a per system basis, it is now possible to blacklist certain talkgroup IDs against ingestion.
- Groups and tags are now defined in their own section, then linked to talkgroup definitions.
- Server logs are now stored in the database and accessed through the administrative dashboard, in addition to the standard output.
- Talkgroups CSV files can now be loaded in from the administrative dashboard.
- Server configuration can be exported/imported to/from a JSON file.
- The downstream id_as property is gone due to its complexity of implementation with the new systems/talkgroups selection dialog for access codes, downstreams and apikeys.
- The keyboard shortcuts are a thing of the past. They caused conflicts with other features.
- Minor changes to the webapp look, less rounded.
- Talkgroup buttons label now wraps on 2 lines.

_v5.1.1_

- Fix database migration script to version 5.1 to filter out duplicate property values on unique fields.
- Fix payload too large error message when saving configuration from the administrative dashboard.
- Bring back the load-rrdb, load-tr and random uuid command line tools.

_v5.1.2_

- Fix config class not returning proper id properties when new records are added.
- Fix database migration script to version 5.1 when on mysql.
- Fix bad logic in apiKey validation.
- Remove the autoJsonMap from the sequelize dialectOptions.
- Client updated to angular 12.

## Version 5.0

- Add rdioScanner.options.autoPopulate which by default is true. The configuration file will now be automatically populated from new received calls with unknown system/talkgroup.
- Add rdioScanner.options.sortTalkgroupswhich by default is false. Sort talkgroups based on their ID.
- Remove default rdioScanner.systems for new installation, since now we have autoPopulate.
- Node modules update.

_v5.0.1_

- Remove the EBU R128 loudness normalization as it's not working as intended.
- Fix the API key validation when using the complex syntax.

_v5.0.2_

- Fix rdioScanner.options.disableAudioConversion which was ignored when true.

_v5.0.3_

- Fix error with docker builds where sequelize can't find the sqlite database.

_v5.0.4_

- Improvement to load-rrdb and load-rr functions.
- Sort groups on the selection panel.
- Allow downstream to other instances running with self-signed certificates.
- Node modules update.

_v5.0.5_

- Node modules security update.
- Improve documentation in regards to minimal Node.js LTS version.
- Add python to build requirements (to be able to build SQLite node module).

## Version 4.9

- Add basic duplicate call detection and rejection.
- Add keyboard shortcuts for the main buttons.
- Add an avoid indicator when the talkgroup is avoided.
- Add an no link indicator when websocket connection is down.
- Node modules update.

_v4.9.1_

- Add EBU R128 loudness normalization.
- dirWatch.type="trunk-recorder" now deletes the JSON file in case the audio file is missing.
- Fix downstream sending wrong talkgroup id.

_v4.9.2_

- Add Config.options.disableKeyboardShortcuts to make everyone a happy camper.

## Version 4.8

- Add downstream.system.id_as property to allow export system with a different id.
- Add system.order for system list ordering on the client side.
- Fix client main screen unscrollable overflow while in landscape.
- Fix issue 26 - date in documentation for mask isn't clear.
- The skip button now also allows you to skip the one second delay between calls.
- Node modules update.

_v4.8.1_

- Refactor panels' back button and make them fixed at the viewport top.
- Node modules update.

_v4.8.2_

- Fix dirWatch.type='sdr-trunk' metatag artist as source is now optional.
- Fix dirWatch.type='sdr-trunk' metatag title as talkgroup.id.
- Web app now running with Angular 11.
- Node modules update.

_v4.8.3_

- Add the ability to overwrite the default dirWatch extension for type sdr-trunk and trunk-recorder.
- Fix dirWatch.disabled being ignored.
- Node modules update.

_v4.8.4_

- Fix the timezone issue when on mariadb.
- Fix downstream sending wrong talkgroup id.
- Node modules security update.

_v4.8.5_

- Fix broken dirwatch.delay.
- Node modules update.

## Version 4.7

- New dirWatch.type='sdr-trunk'.
- New search panel layout with new group and tag filters.
- Add load-tr to load Trunk Recorder talkgroups csv.
- Remove Config.options.allowDownloads, but the feature remains.
- Remove Config.options.useGroup, but the feature remains.
- Bug fixes.

_v4.7.1_

- Fix crash on client when access to talkgroups is restricted with a password.

_v4.7.2_

- Fix Keypad beeps not working on iOS.
- Fix pause not going off due to the above bug.

_v4.7.3_

- Fix websocket not connection on ssl.

_v4.7.4_

- Fix display width too wide when long talkgroup name.

_v4.7.5_

- Fix playback mode getting mixed up if clicking too fast on play.
- Fix side panels background color inheritance.
- Node modules update.

_v4.7.6_

- Fix search results not going back to page 1 when search filters are modified.
- Skip next button no longer emit a denied beeps sequence when pushed while there's no audio playing.
- Node modules update.

## Version 4.6

- Fix documentation in regards to load-rrd in install-github.md.
- Fix database absolute path in config.json.
- Remove config.options.useLed.
- Rename Config.options.keyBeep to Config.options.keypadBeeps.
- Config.options.keypadBeeps now with presets instead of full pattern declaration.
- Bug fixes.

## Version 4.5

- Config.options.keyBeep which by default is true.
- Bug fixes.

## Version 4.4

- Config.systems.talkgroups.patches to group many talkgroups (patches) into one talkgroup.id.
- Config.options now groups allowDownloads, disableAudioConversion, pruneDays, useDimmer, useGroup and useLed options instead of having them spread all over the config file.
- Client will always display talkgroup id on the right side instead of 0 when call is analog.
- Fix annoying bug when next call queued to play is still played even though offline continuous play mode is turned off.
- Talkgroup ID is displayed no matter what and unit ID is displayed only if known.

## Version 4.3

- Add metatags to converted audio files.
- Automatic database migration on startup.
- Client now on Angular 10 in strict mode.
- Dockerized.
- Fix downstream not being triggered when a new call imported.
- Fix dirWatch mask parser and new mask metatags.
- Fix stop button on the search panel when in offline play mode.
- Fix SSL certificate handling.
- Rewritten documentation.

## Version 4.2

- Fix possible race conditions....
- Added websocket keepalive which helps mobile clients when switching from/to wifi/wan.
- Better playback offline mode animations and queue count.
- New dirWatch.mask option to simplify meta data import.

## Version 4.1

- New playback mode.

## Version 4.0

- GraphQL replaced by a pure websocket command and control system.
- `server/.env` replaced by a `server/config.json`.
- Systems are now configured through `server/config.json`, which also invalidate the script `upload-system`.
- Indexes which result in much faster access to archived audio files.
- Add SSL mode.
- Restrict systems/talkgroups access with passwords.
- Directory watch and automatic audio files ingestion.
- Automatic m4a/aac file conversion for better compatibility/performance.
- Selectively share systems/talkgroups to other instances via downstreams.
- Customizable LED colors by systems/talkgroups.
- Dimmable display based on active call.

### Upgrading from version 3

- Your `server/.env` file will be used to create the new `server/config.json` file. Then the `server/.env` will be deleted.
- The `rdioScannerSystems` table will be used to create the _rdioScanner.systems_ within `server/config.json`. Then the `rdioScannerSystems` table will be purged.
- The `rdioScannerCalls` table will be rebuilt, which can be pretty long on some systems.
- It is no longer possible to upload neither your TALKGROUP.CSV nor you ALIAS.CSV files to _Rdio Scanner_. Instead, you have to define them in the `server/config.json` file.

> YOU SHOULD BACKUP YOUR `SERVER/.ENV` FILE AND YOUR DATABASE PRIOR TO UPGRADING, JUST IN CASE. WE'VE TESTED THE UPGRADE PROCESS MANY TIMES, BUT WE CAN'T KNOW FOR SURE IF IT'S GOING TO WORK WELL ON YOUR SIDE.

## Version 3.1

- Client now on Angular 9.
- Display listeners count on the server's end.

## Version 3.0

- Unit aliases support, display names instead of unit ID.
- Download calls from the search panel.
- New configuration options: _allowDownload_ and _useGroup_.

> Note that you can only update from version 2.0 and above. You have to do a fresh install if your actual version is prior to version 2.0.

## Version 2.5

- New group toggle on the select panel.

## Version 2.1

- Various speed improvements for searching stored calls.

## Version 2.0

- Ditched meteor in favour of GraphQL.

## Version 1.0

- First public version.

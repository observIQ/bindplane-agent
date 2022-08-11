# observIQ Distro for OpenTelemetry Updater

The updater is a separate binary that runs as a separate process to update collector artifacts (including the collector itself) when managed by [BindPlane OP](https://github.com/observIQ/bindplane-op).

Because the updater interacts with the service manager, and may edit privileged files, it needs elevated privileges to run (root on Linux + macOS, administrative privileges on Windows).

## Updating Overview

1. The collector receives PackagesAvailable message.
2. The collector downloads tarball containing new updater and updated artifacts (e.g. new collector, plugins, etc.) based on the contents of the PackagesAvailable message.
3. The collector unpacks tarball into `$INSTALL_DIR/tmp/latest`.
4. The collector copies the newest updater binary from `$INSTALL_DIR/tmp/latest` to the working directory.
5. The collector starts the updater in as a separate process in a new process group.
   * If the updater fails to stop the collector within 30 seconds, the collector will kill the updater and abort the update. 

6. The updater starts, then shuts down the collector through the service manager, orphaning the updater process.
7. The updater creates a backup of the current installation directory in `$INSTALL_DIR/tmp/rollback`.
   * If backing up fails for some reason, the updater starts the collector again and exits.
8. The updater installs new artifacts, copying the new files into the the installation directory.
   * If installation fails for some reason, a rollback is initiated.
9. The updater updates the service configuration.
10. The updater starts the collector again, monitoring for collector to be healthy.
    * If the collector is determined to be healthy, the updater exits
    * If the collector is determined unhealthy or doesn't report healthy within 10 seconds, a rollback is initiated. 
11. Upon exit, the updater removes the tmp directory.

## Collector Status Monitoring
The collector saves its current state (installing, installation failed, or installation successful) to a JSON file (`package_statuses.json`) on disk. The updater continuously polls this file for changes in order to detect whether the collector is healthy or not. 

If the collector starts, is able to run successfully, and has its expected version (the version we are upgrading to), the collector will write that the installation was successful to the JSON file. Otherwise, it will write that the installation failed, and will expect the updater to perform a rollback.

If the file indicates the installation failed, or the file still indicates the collector is in an installing state after 10 seconds, a rollback to the previous version is initiated by the updater.

If the file indicates the installation was successful, then the updater exits.

## Updater Rollback
While installing, the updater records a list of actions take (files copied, service actions taken). If something goes wrong during installation, or while monitoring for collector health, then a rollback is initiated.

The rollback will perform the reverse of each action in reverse order. For instance, if the collector binary is replaced with the new collector, then the OpenTelemetry JMX jar is replaced with a new jar, the rollback would copy the backup jar to its previous location, then it would copy the backup collector to its previous location.

Ultimately, this means the rollback process will put the system back into its original state.

After rollback is complete, the updater process exits.

## Debugging
In the event of update failure, you may look in the following places for more information:

* The collector produces logs in the log directory (by default, `$INSTALL_DIR/log/collector.log`).
* The updater produces logs in the collector's log directory, and is always set to log at a debug level (`$INSTALL_DIR/log/updater.log`).
* If the update was initiated through the BindPlane OP UI, an error message should be displayed in the UI.

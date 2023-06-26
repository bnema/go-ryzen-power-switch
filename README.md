# GoRyzenPowerSwitch

GoRyzenPowerSwitch is a small utility that automatically switches power modes on Ryzen-powered laptops based on the current power source. It uses the `ryzenadj` and `cpupower` tools to adjust power and CPU frequency settings.

## Overview

The utility runs in a loop, checking every 10 seconds whether the laptop is running on battery or AC power. Depending on the power source, it applies one of the two following configurations:

- Battery mode: Runs `ryzenadj --power-saving` and `cpupower frequency-set -g conservative`.
- AC mode: Runs `ryzenadj --max-performance` and `cpupower frequency-set -g performance`.

Additionally, it provides two command-line flags `--battery` and `--plugged` to force battery or AC mode, respectively, which can be useful for testing purposes.

## Requirements

- Go 1.20 or later.
- The `ryzenadj` and `cpupower` utilities must be installed and available in your PATH. These tools are used to adjust power and CPU frequency settings. 
- This utility uses the `/sys/class/power_supply/ACAD/online` file to determine whether the laptop is running on AC power. The exact path might vary depending on your distribution and hardware.
- This utility must be run with root privileges because `ryzenadj` and `cpupower` require them to adjust power settings. It will prompt for your sudo password when it starts.
- Developed and tested on Arch Linux. Might work on other distributions, but this is not guaranteed.

## Service Script

If you want to run GoRyzenPowerSwitch permanently, you can create a systemd service:

```bash
sudo nano /etc/systemd/system/goryzenpowerswitch.service
[Unit]
Description=GoRyzenPowerSwitch Service

[Service]
ExecStart=/path/to/goryzenpowerswitch
Restart=on-failure

[Install]
WantedBy=multi-user.target
```
    
Then, enable and start the service:
    
```bash
sudo systemctl enable goryzenpowerswitch
sudo systemctl start goryzenpowerswitch
```


## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
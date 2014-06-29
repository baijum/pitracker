import os
import configparser


def get_config(section="general"):
    """Return test configuration as a dictionary

    Look at `config.ini` for the configuration.
    """
    cwd = os.path.dirname(__file__)
    config_file = os.path.join(cwd, "config.ini")
    cp = configparser.SafeConfigParser()
    cp.read(config_file)
    config = dict(cp.items(section))
    for key, value in config.items():
        env_key = "SEL_" + section.upper() + "_" + key.upper()
        config[key] = os.environ.get(env_key, value)
    return config

general_config = get_config("general")

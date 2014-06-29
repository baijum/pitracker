import time
import subprocess

import pytest

from selenium import webdriver
from config import general_config
import docker


class ResourceHandler(object):

    def __init__(self):
        self.driver = webdriver.Firefox()
        self.driver.maximize_window()
        self.driver.get(general_config['base_url'])
        self.docker_client = docker.Client(base_url='unix://var/run/docker.sock',
                                           version=general_config['docker_version'],
                                           timeout=10)
        self.container = self.docker_client.create_container(general_config['postgresql_image'])
        self.docker_client.start(self.container['Id'], port_bindings={5432: 5432})
        time.sleep(5)
        subprocess.call(general_config['init_docker'], shell=True)


    def release(self):
        self.driver.close()
        self.driver.quit()
        # NOTE: Change to 'stop' after closing db connection using a REST API call
        self.docker_client.kill(self.container['Id'])


def _release_resource_handler(handler):
    """teardown resource_handler"""
    handler.release()


def _get_resource_handler():
    """Factory for resource_handler"""
    resource_handler = ResourceHandler()
    return resource_handler


@pytest.fixture
def resource_handler(request):
    """Create resource_handler funcarg"""
    return request.cached_setup(
        setup=_get_resource_handler,
        teardown=_release_resource_handler,
        scope='function')

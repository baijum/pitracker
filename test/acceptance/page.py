from locators import MainPageLocators
from locators import AddNewProjectPageLocators
from locators import ProjectEditPageLocators
from element import BasePageElement
from selenium.webdriver.support.wait import WebDriverWait
from selenium.common.exceptions import NoSuchElementException
from selenium.webdriver.common.by import By


class BasePage(object):
    """Base page object"""

    implicit_wait = 10

    def __init__(self, resource_handler):
        self.driver = resource_handler.driver
        self.resource_handler = resource_handler

    def wait_until(self, by, value, timeout=5, parent_element=None):
        """Wait until an element is found by a specified value

        This will be useful when part of page is loded though Ajax

        :param value: text to be identified in the page

        :param by: type of the element like available in ``selenium.webdriver.common.by.By`` class

        :param timeout: Time out value, if element is not found within the time out,
            time out error is raised

        :param parent_element: parent element, any object returned by find_elemement_by* methods

        - ID = "id"
        - XPATH = "xpath"
        - LINK_TEXT = "link text"
        - PARTIAL_LINK_TEXT = "partial link text"
        - NAME = "name"
        - TAG_NAME = "tag name"
        - CLASS_NAME = "class name"
        - CSS_SELECTOR = "css selector"

        eg::

          from selenium.webdriver.common.by import By

          wait_until(By.ID, "some_id")
        """
        if parent_element:
            WebDriverWait(parent_element, timeout).until(
                lambda parent_element: self.is_element_present(by, value, parent_element=parent_element))
        else:
            WebDriverWait(self.driver, timeout).until(
                lambda driver: self.is_element_present(by, value))

    def is_element_present(self, by, value, parent_element=None):
        self.driver.implicitly_wait(0)
        if parent_element:
            element_or_driver = parent_element
        else:
            element_or_driver = self.driver
        try:
            element_or_driver.find_element(by, value)
            return True
        except NoSuchElementException:
            return False
        finally:
            # set back to where you once belonged
            self.driver.implicitly_wait(self.implicit_wait)


class NameElement(BasePageElement):

    locator = 'name'


class DescriptionElement(BasePageElement):

    locator = 'description'


class MainPage(BasePage):

    def click_new_project_link(self):
        element = self.driver.find_element(*MainPageLocators.NEW_PROJECT_LINK)
        element.click()

    def click_all_projects_link(self):
        element = self.driver.find_element(*MainPageLocators.ALL_PROJECTS_LINK)
        element.click()


class NewProjectPage(BasePage):

    name_element = NameElement()
    description_element = DescriptionElement()

    def click_create_project_button(self):
        element = self.driver.find_element(*AddNewProjectPageLocators.CREATE_PROJECT_BUTTON)
        element.click()


class AllProjectsPage(BasePage):

    def is_project_exists(self, project_name):
        return self.is_element_present(By.LINK_TEXT, project_name)

    def click_project_link(self, project_name):
        element = self.driver.find_element_by_link_text(project_name)
        element.click()


class ProjectViewPage(BasePage):

    def click_edit_project_link(self):
        element = self.driver.find_element_by_link_text('Edit')
        element.click()


class ProjectEditPage(BasePage):

    name_element = NameElement()
    description_element = DescriptionElement()

    def click_update_project_button(self):
        element = self.driver.find_element(*ProjectEditPageLocators.UPDATE_PROJECT_BUTTON)
        element.click()

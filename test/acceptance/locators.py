from selenium.webdriver.common.by import By


class MainPageLocators(object):
    NEW_PROJECT_LINK = (By.LINK_TEXT, 'New Project')
    ALL_PROJECTS_LINK = (By.LINK_TEXT, 'All Projects')


class AddNewProjectPageLocators(object):
    CREATE_PROJECT_BUTTON = (By.XPATH, '//button[text()="Create Project"]')


class AllProjectsPageLocators(object):
    pass


class ProjectPage(object):
    CREATE_PROJECT_BUTTON = (By.XPATH, '//a[@href="#/projects/1/edit"]')


class ProjectEditPageLocators(object):
    UPDATE_PROJECT_BUTTON = (By.XPATH, '//button[text()="Update Project"]')

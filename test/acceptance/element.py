from selenium.webdriver.support.ui import WebDriverWait


class BasePageElement(object):

    def __set__(self, obj, value):
        # fetch driver from caller
        driver = obj.driver
        # set element
        WebDriverWait(driver, 100).until(
            lambda driver: driver.find_element_by_id(self.locator))
        # force focus on element before send_keys
        driver.find_element_by_id(self.locator).click()
        if value.startswith(':no-clear:'):
            value = value.replace(':no-clear:', '')
        else:
            driver.find_element_by_id(self.locator).clear()
        driver.find_element_by_id(self.locator).send_keys(value)

    def __get__(self, obj, owner):
        try:
            driver = obj.driver
        except AttributeError:
            return
        WebDriverWait(driver, 100).until(
            lambda driver: driver.find_element_by_id(self.locator))
        element = driver.find_element_by_id(self.locator)
        return element.get_attribute("value")


class DropDownElement(object):

    def __set__(self, obj, value):
        driver = obj.driver
        WebDriverWait(driver, 100).until(
            lambda driver: driver.find_element_by_xpath("//select[@name='%s']/option" % self.locator))
        if value.startswith(':value:'):
            value = value.replace(':value:', '')
            select = driver.find_element_by_xpath("//select[@name='%s']/option[normalize-space(@value)='%s']" % (self.locator, value))
        elif value.startswith(':partial-text:'):
            value = value.replace(':partial-text:', '')
            select = driver.find_element_by_xpath("//select[@name='%s']/option[contains(normalize-space(text()),'%s')]" % (self.locator, value))
        else:
            select = driver.find_element_by_xpath("//select[@name='%s']//option[normalize-space(text())='%s']" % (self.locator, value))
        select.click()

    def __get__(self, obj, owner):
        try:
            driver = obj.driver
        except AttributeError:
            return
        element = WebDriverWait(driver, 100).until(
            lambda driver: driver.find_element_by_name(self.locator))
        element = driver.find_element_by_name(self.locator)
        element_value = element.get_attribute("value")
        all_options = element.find_elements_by_tag_name("option")
        for options in all_options:
            if options.get_attribute("value") == element_value:
                return str(options.text)

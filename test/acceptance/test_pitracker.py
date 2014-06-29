import page
import pytest


class TestPITracker:

    @pytest.mark.parametrize(("project_name",),
                             [("test1",),])
    def test_project_creation(self, resource_handler, project_name):
        main_page = page.MainPage(resource_handler)
        main_page.click_new_project_link()
        new_project_page = page.NewProjectPage(resource_handler)
        new_project_page.name_element = project_name
        new_project_page.description_element = project_name
        new_project_page.click_create_project_button()
        main_page.click_all_projects_link()
        all_projects_page = page.AllProjectsPage(resource_handler)
        assert all_projects_page.is_project_exists(project_name), "Cannot find project: %s" % project_name
        all_projects_page.click_project_link(project_name)
        project_view_page = page.ProjectViewPage(resource_handler)
        project_view_page.click_edit_project_link()
        project_edit_page = page.ProjectEditPage(resource_handler)
        project_edit_page.name_element = "new " + project_name
        project_edit_page.description_element = "new " + project_name
        project_edit_page.click_update_project_button()
        main_page = page.MainPage(resource_handler)
        main_page.click_all_projects_link()
        all_projects_page = page.AllProjectsPage(resource_handler)
        assert all_projects_page.is_project_exists("new " + project_name), "Cannot find project: %s" % "new " + project_name

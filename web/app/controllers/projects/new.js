import Ember from "ember";

export default Ember.Controller.extend({
  actions: {
    createProject: function() {
      var self = this;
      var project = self.store.createRecord('project', {
        name: self.get('name'),
        _id: self.get('_id'),
        description: self.get('description')
      });
      project.save().then(function() {
	self.transitionTo('index');
      }, function(response) {
	if (response.status === 409) {
	  console.log("Project already exists");
	}
      });
    }
  }
});

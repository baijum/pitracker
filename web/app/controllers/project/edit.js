import Ember from "ember";

export default Ember.Controller.extend({
  actions: {
    updateProject: function(project) {
      var self = this;
      project.save();
      this.transitionToRoute('index');
    }
  }
});

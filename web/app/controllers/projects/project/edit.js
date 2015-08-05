import Ember from 'ember';

export default Ember.Controller.extend({
  actions: {
    save: function(project) {
      project.save();
      this.transitionTo('index');
    }
  }
});

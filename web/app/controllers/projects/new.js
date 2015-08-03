import Ember from 'ember';

export default Ember.Controller.extend({
  actions: {
    save: function() {
      var self = this;
      var name = self.get('name');
      var description = self.get('description');
      var project = self.store.createRecord('project', {
        name: name,
        description: description
      });
      project.save();
      self.transitionToRoute('index');
    }
  }
});

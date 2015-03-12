import Ember from "ember";
 
export default Ember.Controller.extend({
  actions: {
    createProject: function() {
      console.log(this.name);
      console.log(this.description);

      var project = this.store.createRecord('project', {
        name: this.name,
        description: this.description
      });
 
      project.save();
    }
  }
});

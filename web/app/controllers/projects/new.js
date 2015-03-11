import Ember from "ember";
 
export default Ember.Controller.extend({
  actions: {
    createProject: function() {
      console.log(this.name);
      console.log(this.description);
    }
  }
});

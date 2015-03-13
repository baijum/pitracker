import Ember from "ember";

export default Ember.Route.extend({
  actions: {
    error: function(error, transition) {
      // handle the error
      console.log(error.message);
      console.log(transition);
    }
  }
});

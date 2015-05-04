import Ember from "ember";

export default Ember.Controller.extend({
  actions: {
    updateItem: function(item) {
      var self = this;
      item.save();
      self.transitionToRoute('index');
    }
  }
});

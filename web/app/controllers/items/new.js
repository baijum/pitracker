import Ember from 'ember';

export default Ember.Controller.extend({
  actions: {
    createItem: function() {
      var self = this;
      var item = self.store.createRecord('item', {
	title: self.get('title'),
	description: self.get('description')
      });
      item.save().then(function() {
	self.transitionTo('index');
      }, function(response) {
	if (response.status === 409) {
	  console.log("Item already exists");
	}
      });
    }
  }
});

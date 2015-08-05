import Ember from 'ember';
import AuthenticatedRouteMixin from 'simple-auth/mixins/authenticated-route-mixin';

export default Ember.Controller.extend(AuthenticatedRouteMixin, {
  actions: {
    save: function() {
      var project = this.get('project');
      var title = this.get('title');
      var description = this.get('description');
      var item = this.store.createRecord('item', {
	project: project,
        title: title,
        description: description
      });
      item.save();
      this.transitionToRoute('index');
    }
  }
});

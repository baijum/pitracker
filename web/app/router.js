import Ember from 'ember';
import config from './config/environment';

var Router = Ember.Router.extend({
  location: config.locationType
});

Router.map(function() {
  this.route('login');
  this.route('projects', function() {
    this.route('new');
    this.route('project', { path: ':project_id' }, function() {
      this.route('edit');
    });
  });
});

export default Router;

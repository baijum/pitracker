import Ember from 'ember';
import config from './config/environment';

var Router = Ember.Router.extend({
  location: config.locationType
});

Router.map(function() {
  this.route('about');
  this.resource('projects', function() {
    this.route('new');
    this.resource('project', { path: ':project_id' }, function() {
      this.route('edit');
      this.route('archive');
    });
    this.resource('milestones', function() {
      this.route('new');
      this.resource('milestone', { path: ':milestone_id' });
    });
  });
  this.resource('events', function() {
    this.route('view');
  });
});

export default Router;

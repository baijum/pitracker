import Ember from 'ember';
import config from './config/environment';

var Router = Ember.Router.extend({
  location: config.locationType
});

Router.map(function() {
  this.route('about');
  this.resource('projects', function() {
    this.route('new');
    this.resource('project', { path: ':project_name' }, function() {
      this.route('edit');
      this.route('archive');
    });
  });
  this.resource('items', function() {
    this.route('new');
    this.resource('item', { path: ':item_id' }, function() {
      this.route('edit');
      this.route('archive');
//      this.resource('comment', { path: ':comment_id' }, function() {
//	this.route('edit');
//      });
    });
  });
});

export default Router;

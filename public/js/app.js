App = Ember.Application.create();

// Routes
App.Router.map(function() {
    this.resource('items', function() {
        this.route('new');
	this.resource('item', { path: ':item_id' }, function() {
	    this.route('edit');
	});
    });
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
    this.route('about');
    this.route('login');
    this.route('logout');
});

// Data adapter
App.ApplicationAdapter = DS.RESTAdapter.extend({
    namespace: 'api/v1',
    headers: {
        Authorization: 'Bearer '+ localStorage.getItem('access_token')
    }
});

// Models
App.Project = DS.Model.extend({
    name: DS.attr('string'),
    description: DS.attr('string'),
    items: DS.hasMany('item', {
        async: true
    })
});

App.Item = DS.Model.extend({
    project: DS.belongsTo('project'),
    title: DS.attr('string'),
    description: DS.attr('string')
});

// Route definitions
App.AuthenticatedRoute = Ember.Route.extend({

    beforeModel: function(transition) {
        if (!localStorage.access_token) {
            this.redirectToLogin(transition);
        }
    },

    redirectToLogin: function(transition) {
        var loginController = this.controllerFor('login');
        loginController.set('attemptedTransition', transition);
        this.transitionTo('login');
    },

    getJSONWithToken: function(url) {
        var access_token = this.controllerFor('login').get('access_token');
        return $.getJSON(url, { access_token: access_token });
    },

    actions: {
        error: function(reason, transition) {
            if (reason.status === 401) {
                this.redirectToLogin(transition);
            } else {
                alert(this);
                alert('Something went wrong');
            }
        }
    }

});

// non-authenticated route to login
App.LoginRoute = Ember.Route.extend({
});

App.ProjectsIndexRoute = App.AuthenticatedRoute.extend({
    model: function() {
        return this.store.find('project');
    }
});

App.ProjectsNewRoute = Ember.Route.extend({
});

App.ProjectsNewController = Ember.Controller.extend({
    actions: {
        save: function() {
	    var name = this.get('name');
            var description = this.get('description');
            var project = this.store.createRecord('project', {
                name: name,
                description: description
            });
            project.save();
            this.transitionToRoute('index');
        }
    }
});

App.ProjectRoute = Ember.Route.extend({
    model: function(params) {
	return this.store.find('project', params.project_id);
    }
});

App.IndexController = Ember.Controller.extend({
    logout: function() {
        delete localStorage.access_token;
        this.controllerFor('login').set('loggedIn', false);
        this.transitionTo('index');
    }
});

App.LoginController = Ember.Controller.extend({
    setupController: function(controller, context) {
        controller.reset();
    },

    reset: function() {
        delete localStorage.access_token;
        this.controllerFor('login').set('loggedIn', false);
        this.setProperties({
            username: '',
            password: '',
            errorMessage: '',
            loggedIn: false
        });
    },

    access_token: localStorage.access_token,

    tokenChanged: function() {
        localStorage.access_token = this.get('access_token');
    }.observes('access_token'),

    actions: {
        login: function() {

            var self = this, data = this.getProperties('username', 'password');

            // Clear out any error messages.
            this.set('errorMessage', null);

            $.post('/api/v1/auth', data).then(function(response) {

                self.set('errorMessage', response.message);
                if (response.success) {

                    self.set('access_token', response.access_token);
                    self.set('loggedIn', true);

                    // FIXME: Is there any better way to handle this ?
                    App.ApplicationAdapter.reopen({
                        headers: {Authorization: 'Bearer '+ localStorage.access_token }
                    });

                    var attemptedTransition = self.get('attemptedTransition');
                    var previousTransition = self.get('previousTransition');

                    if (attemptedTransition) {
                        alert('Hello 4');
                        // FIXME: This seems to be not working
                        self.attemptedTransition.retry();
                        //self.set('attemptedTransition', null);
                        //self.set('attemptedTransition', null);
                        //self.transitionToRoute('index');
                        alert('Hello 6');
                    } else {
                        // Redirect to 'index' by default.
                        self.transitionToRoute('index');
                        alert('Hello 5');
                    }
                }
            });

        }
    }
});

App.ItemsIndexRoute = Ember.Route.extend({
    model: function() {
        return this.store.find('item');
    }
});

App.ItemsNewRoute = Ember.Route.extend({
});

App.ItemsNewController = Ember.Controller.extend({
    actions: {
	save: function() {
	    var title = this.get('title');
            var description = this.get('description');
            var item = this.store.createRecord('item', {
                title: title,
                description: description
            });
            item.save();
            this.transitionToRoute('index');
	}
    }
});

App.ItemRoute = Ember.Route.extend({
    model: function(params) {
        return this.store.find('item', params.item_id);
    }
});

App.ItemIndexRoute = Ember.Route.extend({
});

App.ItemEditRoute = Ember.Route.extend({
    actions: {
        save: function(item) {
            item.save();
            this.transitionTo('index');
        }
    }
});

App.ProjectEditRoute = Ember.Route.extend({
    actions: {
        save: function(project) {
            project.save();
            this.transitionTo('index');
        }
    }
});

App.ProjectArchiveRoute = Ember.Route.extend({
    actions: {
        save: function(project) {
            project.deleteRecord();
            project.save();
            this.transitionTo('index');
        }
    }
});

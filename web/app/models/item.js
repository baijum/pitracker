import DS from 'ember-data';

export default DS.Model.extend({
    project: DS.belongsTo('project'),
    title: DS.attr('string'),
    description: DS.attr('string'),
    comments: DS.hasMany('comment', {
      async: true
    })
});

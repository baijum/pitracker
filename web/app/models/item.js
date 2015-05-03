import DS from 'ember-data';

export default DS.Model.extend({
//  _id: DS.attr('string'),
  title: DS.attr('string'),
  description: DS.attr('string'),
  comments: DS.hasMany('comment')
});

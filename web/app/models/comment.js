import DS from 'ember-data';

export default DS.Model.extend({
  comment: DS.attr('string'),
  item: DS.belongsTo('item')
});

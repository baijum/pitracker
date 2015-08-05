import DS from 'ember-data';

export default DS.Model.extend({
    item: DS.belongsTo('item'),
    content: DS.attr('string')
});

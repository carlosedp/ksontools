package yaml2jsonnet

type Properties map[string]interface{}

func matchProp(root *Node, t map[interface{}]interface{}, name, path string) error {
	return nil

	// node, err := root.Property(name)
	// if err != nil {
	// 	return errors.Wrapf(err, "inspect property %s", name)
	// }

	// if node.IsMixin {
	// 	logrus.WithField("mixinName", node.name).Info("found mixin")

	// }

	// for k, v := range t {
	// 	k1 := k.(string)
	// 	setter, err := node.FindFunction(name, k1)
	// 	if err != nil {
	// 		logger.Warnf("%s is a mixin", k1)
	// 		continue
	// 	}

	// 	if err := comp.AddParam(k1, v); err != nil {
	// 		return errors.Wrap(err, "add param")
	// 	}

	// 	builders = append(builders, fmt.Sprintf("%s(%s)", setter, k1))
	// }

	// if node.IsMixin && len(builders) > 0 {
	// 	method := fmt.Sprintf("%s.mixin.%s.%s",
	// 		d.GVK.Kind,
	// 		node.name,
	// 		strings.Join(builders, "."))

	// 	val := NewDeclarationApply(method)

	// 	mixinName := fmt.Sprintf("%s%s", d.GVK.Kind, strings.Title(node.name))

	// 	decl := Declaration{
	// 		Name:  mixinName,
	// 		Value: val,
	// 	}
	// 	comp.AddDeclaration(decl)

	// 	mixinNames = append(mixinNames, mixinName)
	// }
}

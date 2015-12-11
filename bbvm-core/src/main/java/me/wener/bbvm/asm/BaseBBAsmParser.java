package me.wener.bbvm.asm;

import com.google.common.collect.Maps;
import me.wener.bbvm.vm.Operand;

import java.util.Map;

/**
 * @author wener
 * @since 15/12/10
 */
public class BaseBBAsmParser {

    Map<String, Label> labels = Maps.newHashMap();

    static String labelName(Token token) {
        String image = token.image;
        if (image.endsWith(":")) {
            image = image.substring(image.length() - 1);
        }
        return image;
    }

    void jjtreeOpenNodeScope(Node n) {
    }

    void jjtreeCloseNodeScope(Node n) {
    }

    public void addLabel(Token token) {
        String name = labelName(token);
        Label label = labels.get(name);
        if (label != null) {
            throw new RuntimeException("Label already exists " + label);
        }
        label = new Label(name, token);
        labels.put(label.name, label);
    }

    public void addLabelOperand(Token token, Operand operand) {
        operand.setValue(-1);
        String name = labelName(token);
        Label label = labels.get(name);
        if (label == null) {
            label = new Label(name);
            labels.put(label.name, label);
        }
        label.addOperand(token, operand);
    }

    public void checkLabel() {
        // All label are addressed
        for (Label label : labels.values()) {
            if (label.token == null) {
                throw new RuntimeException("Undefined label " + label);
            }
            if (label.address < 0) {
                throw new RuntimeException("Undressed label " + label);
            }
        }
    }
}

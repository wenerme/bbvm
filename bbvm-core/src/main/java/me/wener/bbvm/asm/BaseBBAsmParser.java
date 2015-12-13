package me.wener.bbvm.asm;

import com.google.common.collect.Lists;
import com.google.common.collect.Maps;
import io.netty.buffer.ByteBuf;
import me.wener.bbvm.vm.Instruction;
import me.wener.bbvm.vm.Operand;
import me.wener.bbvm.vm.SymbolTable;
import me.wener.bbvm.vm.Symbols;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.LinkedList;
import java.util.List;
import java.util.Map;

/**
 * @author wener
 * @since 15/12/10
 */
class BaseBBAsmParser {
    protected final static Logger log = LoggerFactory.getLogger(BBAsmParser.class);
    LinkedList<Assembly> assemblies = Lists.newLinkedList();
    Map<String, Label> labels = Maps.newHashMap();

    static String labelName(Token token) {
        String image = token.image;
        if (image.endsWith(":")) {
            image = image.substring(0, image.length() - 1).trim();
        }
        return image;
    }

    void jjtreeOpenNodeScope(Node n) {
    }

    void jjtreeCloseNodeScope(Node n) {
    }

    public SymbolTable createSymbolTable() {
        return Symbols.table(labels.values());
    }

    public Label addLabel(Token token) {
        String name = labelName(token);
        Label label = labels.get(name);
        if (label != null) {
            if (label.getToken() != null)
                throw new RuntimeException("Label already exists " + label);
            label.setToken(token);
        } else {
            label = new Label(name, token);
        }
        labels.put(label.name, label);
        assemblies.add(label);
        return label;
    }

    public void add(Instruction instruction) {
        assemblies.add(new Inst(instruction));
    }

    public void add(Assembly assembly) {
        if (assembly instanceof Label) {
            Label label = (Label) assembly;
            Label old = labels.get(label.getName());

            if (old.getToken() == null) {
                label.operands.addAll(old.operands);
                labels.put(label.name, label);
            } else {
                throw new RuntimeException(String.format("Detect conflict label %s %s,%s <> %s,%s"
                        , label.getName(), label.token.beginLine, label.token.beginColumn, old.token.beginLine, old.token.beginColumn));
            }
        }
        assemblies.add(assembly);
    }

    public void addLabelOperand(Token token, Operand operand) {
        operand.setInternal(-1);
        String name = labelName(token);
        Label label = labels.get(name);
        if (label == null) {
            label = new Label(name);
            labels.put(label.name, label);
        }
        label.addOperand(token, operand);
    }

    public List<Assembly> getAssemblies() {
        return assemblies;
    }

    public void checkLabel() {
        // All label are addressed
        for (Label label : labels.values()) {
            if (label.getToken() == null) {
                throw new RuntimeException("Undefined label " + label);
            }
            if (label.value < 0) {
                throw new RuntimeException("Undressed label " + label);
            }
        }
    }

    public ByteBuf write(ByteBuf buf) {
        for (Assembly assembly : assemblies) {
            assembly.write(buf);
        }
        return buf;
    }

    /**
     * @return Estimated length
     */
    public int estimateLabelAddress() {
        int pos = 0;
        for (Assembly assembly : assemblies) {
            if (assembly instanceof Label) {
                ((Label) assembly).setValue(pos);
            } else {
                pos += assembly.length();
            }
        }
        return pos;
    }

    public void addComment(Token token, boolean isFullLine) {
        if (isFullLine) {
            assemblies.add(new Comment(token));
        } else {
            assemblies.getLast().setComment(new Comment(token));
        }
    }
}

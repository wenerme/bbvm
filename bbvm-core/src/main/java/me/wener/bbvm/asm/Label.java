package me.wener.bbvm.asm;

import com.google.common.base.MoreObjects;
import com.google.common.collect.Lists;
import me.wener.bbvm.vm.Operand;
import me.wener.bbvm.vm.Symbol;

import java.util.List;

/**
 * @author wener
 * @since 15/12/11
 */
public class Label extends AbstractAssembly implements Symbol, Assembly {
    String name;
    int value = -1;
    /**
     * Operands referenced to this label
     */
    List<OperandInfo> operands = Lists.newArrayList();
    Token token;

    public Label() {
    }

    public Label(String name, Token token) {
        this.token = token;
        this.name = name;
    }

    public Label(String name) {
        this.name = name;
    }

    public Token getToken() {
        return token;
    }

    public Label setToken(Token token) {
        this.token = token;
        return this;
    }

    public int getValue() {
        return value;
    }

    public Label setValue(int value) {
        this.value = value;
        for (OperandInfo info : operands) {
            info.operand.setValue(value);
        }
        return this;
    }

    @Override
    public String getName() {
        return name;
    }

    @Override
    public String toString() {
        return MoreObjects.toStringHelper(this)
                .add("name", name)
                .add("value", value)
                .add("comment", comment)
                .add("line", token != null ? token.beginLine : -1)
                .add("column", token != null ? token.beginColumn : -1)
                .toString();
    }

    public void addOperand(Token token, Operand operand) {
        operand.setSymbol(this).setValue(value);
        operands.add(new Label.OperandInfo(operand, token));
    }

    @Override
    public Type getType() {
        return Type.LABEL;
    }

    @Override
    public String toAssembly() {
        return token.toString() + commentAssembly();
    }


    static class OperandInfo {
        Operand operand;
        Token token;

        public OperandInfo(Operand operand, Token token) {
            this.operand = operand;
            this.token = token;
        }
    }
}
